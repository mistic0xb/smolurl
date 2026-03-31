import http from 'k6/http';
import { check } from 'k6';

const codes = JSON.parse(open('./codes.json'));

export const options = {
    scenarios: {
        stress: {
            executor: 'ramping-arrival-rate',
            startRate: 10,
            timeUnit: '1s',
            preAllocatedVUs: 3000,
            maxVUs: 5000,
            stages: [
                { duration: '30s', target: 100 },
                { duration: '1m', target: 500 },
                { duration: '1m', target: 1500 },
                { duration: '1m', target: 2000 },
                { duration: '1m', target: 3000 },
                { duration: '30s', target: 0 },
            ],
        },
    },
    thresholds: {
        http_req_failed: ['rate<0.05'],
        http_req_duration: ['p(95)<2000', 'p(99)<3000'],
    },
};

export default function () {
    const code = codes[Math.floor(Math.random() * codes.length)];
    const res = http.get(
        `http://localhost:8080/${code}`,
        {
            timeout: '5s',
            redirects: 0,
        }
    );
    check(res, {
        'status is 3xx redirect': (r) => r.status >= 300 && r.status < 400,
        'has Location header': (r) => r.headers['Location'] !== undefined,
    });
}