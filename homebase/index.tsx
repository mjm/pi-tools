import React from "react";
import ReactDOM from "react-dom";

import {App} from "com_github_mjm_pi_tools/homebase/App";
import "com_github_mjm_pi_tools/homebase/styles/app.css";

// @ts-ignore
ReactDOM.unstable_createRoot(
    document.getElementById("root")
).render(<App/>);
