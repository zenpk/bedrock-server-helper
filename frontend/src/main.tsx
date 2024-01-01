import React from "react";
import ReactDOM from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import "./global.css";
import { authorization, redirectLogin } from "./util/myoauth.ts";
import { STORAGE_ACCESS_TOKEN } from "./util/constants.ts";
import { Home } from "./page/Home.tsx";
import { WorldPage } from "./page/World.tsx";
import { WorldsContextProvider } from "@/util/WorldsContext.tsx";

const router = createBrowserRouter([
  {
    path: "/",
    element: <Home />,
  },
  {
    path: "/worlds/:worldId",
    element: <WorldPage />,
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
      <WorldsContextProvider>
        <RouterProvider router={router} />
      </WorldsContextProvider>
    </React.StrictMode>,
  );
}
