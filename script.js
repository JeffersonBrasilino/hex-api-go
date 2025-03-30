import http from "k6/http";
import { sleep } from "k6";
import { check } from "k6";

export const options = {
  // A number specifying the number of VUs to run concurrently.
  vus: 100,
  //iterations: 300,
  // A string specifying the total duration of the test run.
  duration: "30s",
};

// The function that defines VU logic.
//
// See https://grafana.com/docs/k6/latest/examples/get-started-with-k6/ to learn more
// about authoring k6 scripts.
//
export default function () {
  const res = http.get("http://localhost:3000/users?data-source=gatewayB");
  check(res, {
    "is status 200": (r) => r.status === 200,
  });
  //sleep(1);
}
