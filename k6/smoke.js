import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    scenarios: {
        stress: {
            executor: 'ramping-arrival-rate',
            startRate: 10,
            timeUnit: '1s',
            preAllocatedVUs: 500,
            maxVUs: 5000,
            stages: [
                { duration: '30s', target: 100 },
                { duration: '1m', target: 1000 },
                { duration: '1m', target: 3000 },
                { duration: '1m', target: 5000 },
                { duration: '30s', target: 0 },
            ],
        },
    },
    thresholds: {
        http_req_failed: ['rate<0.05'],
        http_req_duration: ['p(95)<2000', 'p(99)<3000'],  // added p99
    },
};

export default function() {
    const res = http.post(
        'http://localhost:8080/api/v1/url',
        JSON.stringify({
            original_url: `https://example.com/${__VU}-${__ITER}`,
            expiration_time: 30,
        }),
        {
            headers: { 'Content-Type': 'application/json' },
            timeout: '5s',  // don't let hung requests pile up forever
        }
    );

    check(res, {
        'status is 201': (r) => r.status === 201,
        'has smol_url': (r) => r.status === 201 && JSON.parse(r.body).smol_url !== undefined,
    });

    sleep(0.5);
}
