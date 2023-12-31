import React from "react";
import ReactDOM from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import "./global.css";
import { authorization, redirectLogin } from "./util/myoauth.ts";
import { STORAGE_ACCESS_TOKEN } from "./util/constants.ts";
import { Home } from "./page/Home.tsx";

const router = createBrowserRouter([
  {
    path: "/",
    element: <Home />,
  },
  {
    path: "/:worldId",
    element: <></>,
  },
]);

const urlParams = new URLSearchParams(window.location.search);
if (urlParams.get("authorizationCode")) {
  authorization();
} else {
  if (!window.localStorage.getItem(STORAGE_ACCESS_TOKEN)) {
    redirectLogin();
  }
  ReactDOM.createRoot(document.getElementById("root")!).render(
    <React.StrictMode>
      <RouterProvider router={router} />
    </React.StrictMode>,
  );
}
