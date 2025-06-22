import http from "k6/http";
import { sleep } from "k6";
import { check } from "k6";

export const options = {
  scenarios: {
    constant_request_rate: {
      executor: "constant-vus",
      vus: 500,
      duration: "1m",
    },

    /* one_parallel_request: {
      executor: "per-vu-iterations",
      vus: 3,
      iterations: 1,
      gracefulStop: "5s",
    }, */
  },
  thresholds: {
    //http_req_duration: ["p(95)<500"], // 95% das requisições devem ser menores que 500ms
    http_req_failed: ["rate<0.01"], // Menos de 1% das requisições devem falhar
    //http_reqs: ["count>1000"], // Mais de 1000 requisições totais
  },
};

// The function that defines VU logic.
//
// See https://grafana.com/docs/k6/latest/examples/get-started-with-k6/ to learn more
// about authoring k6 scripts.
//
export default function () {
  const payload = JSON.stringify({
    username: "pegasus",
    senha: "123456",
  });
  const res = http.post("http://localhost:3000/users/create", {
    username: "pegasus",
    senha: "123456",
  });
  check(res, {
    "is status 200": (r) => r.status === 200,
  });
  sleep(1);
}
