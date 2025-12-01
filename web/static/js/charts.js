function LongitudeChart() {
    const labels = chartData1.map(d => d.x);
    const data = chartData1.map(d => d.y);

    const ctx = document.getElementById('lineChart').getContext('2d');
    new Chart(ctx, {
        type: 'line',
        data: {
            labels: labels,
            datasets: [{
                label: '',
                data: data,
                borderColor: '#d63384',
                backgroundColor: 'rgba(214, 51, 132, 0.2)',
                tension: 0.3,
                pointStyle: false,
            }]
        },
        options: {
            responsive: true,
            plugins: {
                legend: { display: false }
            },
            scales: {
                y: { beginAtZero: true }
            }
        }
    });
}

function DistrictChart() {
    const labels = chartData2.map(d => d.x);
    const data = chartData2.map(d => d.y);

    const ctx = document.getElementById('barChart').getContext('2d');
    new Chart(ctx, {
        type: 'bar',
        data: {
            labels: labels,
            datasets: [{
                label: '',
                data: data,
                backgroundColor: '#ffb6c1',
                borderRadius: 8
            }]
        },
        options: {
            responsive: true,
            plugins: {
                legend: { display: false }
            },
            scales: {
                y: { beginAtZero: true }
            }
        }
    });
}

document.addEventListener("DOMContentLoaded", () => {
    LongitudeChart()
    DistrictChart()
});
