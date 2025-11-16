package router

import (
	"net/http"
	"strings"
)

// HandlerFunc тот же тип, что и в net/http, но мы оборачиваем для middleware.
type HandlerFunc func(w http.ResponseWriter, r *http.Request, params Params)

// Params хранит параметры пути (например :id).
type Params map[string]string

// MiddlewareFunc — функция middleware.
type MiddlewareFunc func(HandlerFunc) HandlerFunc

// node — узел для простого trie (разделён по сегментам).
type node struct {
	segment    string
	children   []*node
	paramChild *node                  // child for :param
	catchAll   *node                  // child for *wildcard
	handlers   map[string]HandlerFunc // method -> handler
}

// Router структура маршрутизатора.
type Router struct {
	root             *node
	middlewares      []MiddlewareFunc
	notFound         http.HandlerFunc
	methodNotAllowed http.HandlerFunc
}

// New создаёт новый Router.
func New() *Router {
	return &Router{
		root: &node{
			segment:  "/",
			children: []*node{},
			handlers: map[string]HandlerFunc{},
		},
		notFound: func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		},
		methodNotAllowed: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("405 method not allowed"))
		},
	}
}

// Use добавляет middleware (в порядке вызова).
func (rt *Router) Use(m MiddlewareFunc) {
	rt.middlewares = append(rt.middlewares, m)
}

// Handle регистрирует обработчик для метода и пути.
func (rt *Router) Handle(method, path string, h HandlerFunc) {
	if path == "" || path[0] != '/' {
		panic("path must start with '/'")
	}
	segments := splitPath(path)
	cur := rt.root
	for i, seg := range segments {
		isParam := len(seg) > 0 && seg[0] == ':'
		isCatchAll := len(seg) > 0 && seg[0] == '*'

		var next *node

		// catch-all can only be last
		if isCatchAll {
			if i != len(segments)-1 {
				panic("catch-all must be the last segment")
			}
			if cur.catchAll == nil {
				cur.catchAll = &node{segment: seg, handlers: map[string]HandlerFunc{}}
			}
			next = cur.catchAll
		} else if isParam {
			if cur.paramChild == nil {
				cur.paramChild = &node{segment: seg, handlers: map[string]HandlerFunc{}}
			}
			next = cur.paramChild
		} else {
			// static child
			for _, c := range cur.children {
				if c.segment == seg {
					next = c
					break
				}
			}
			if next == nil {
				next = &node{segment: seg, handlers: map[string]HandlerFunc{}}
				cur.children = append(cur.children, next)
			}
		}
		cur = next
	}
	if cur.handlers == nil {
		cur.handlers = map[string]HandlerFunc{}
	}
	cur.handlers[strings.ToUpper(method)] = h
}

// GET/POST helpers
func (rt *Router) GET(path string, h HandlerFunc)  { rt.Handle("GET", path, h) }
func (rt *Router) POST(path string, h HandlerFunc) { rt.Handle("POST", path, h) }

// ServeHTTP делает Router совместимым с net/http.
func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := r.Method
	handler, params := rt.match(method, path)
	if handler == nil {
		// если путь найден, но метод нет — 405, иначе 404
		if rt.pathExists(path) {
			rt.methodNotAllowed(w, r)
			return
		}
		rt.notFound(w, r)
		return
	}
	// применяем middleware в обратном порядке (wrap)
	final := handler
	for i := len(rt.middlewares) - 1; i >= 0; i-- {
		final = rt.middlewares[i](final)
	}
	final(w, r, params)
}

// match ищет обработчик и собирает параметры для конкретного метода.
func (rt *Router) match(method, path string) (HandlerFunc, Params) {
	segments := splitPath(path)
	cur := rt.root
	params := Params{}

	for i, seg := range segments {
		// try static children first
		var matched *node
		for _, c := range cur.children {
			if c.segment == seg {
				matched = c
				break
			}
		}
		if matched != nil {
			cur = matched
			continue
		}
		// param child?
		if cur.paramChild != nil {
			paramName := strings.TrimPrefix(cur.paramChild.segment, ":")
			params[paramName] = seg
			cur = cur.paramChild
			continue
		}
		// catch-all?
		if cur.catchAll != nil {
			paramName := strings.TrimPrefix(cur.catchAll.segment, "*")
			rest := strings.Join(segments[i:], "/")
			params[paramName] = rest
			cur = cur.catchAll
			break
		}
		// nothing matched
		return nil, nil
	}
	// found node — look up handler by method
	if cur.handlers == nil {
		return nil, nil
	}
	if handler, ok := cur.handlers[strings.ToUpper(method)]; ok {
		return handler, params
	}
	return nil, nil
}

// pathExists проверяет, есть ли путь вообще (без учёта метода).
func (rt *Router) pathExists(path string) bool {
	segments := splitPath(path)
	cur := rt.root
	for i, seg := range segments {
		var matched *node
		for _, c := range cur.children {
			if c.segment == seg {
				matched = c
				break
			}
		}
		if matched != nil {
			cur = matched
			continue
		}
		if cur.paramChild != nil {
			cur = cur.paramChild
			continue
		}
		if cur.catchAll != nil {
			cur = cur.catchAll
			break
		}
		_ = i
		return false
	}
	return cur != nil && len(cur.handlers) > 0
}

func splitPath(p string) []string {
	if p == "/" || p == "" {
		return []string{}
	}
	p = strings.Trim(p, "/")
	if p == "" {
		return []string{}
	}
	return strings.Split(p, "/")
}
