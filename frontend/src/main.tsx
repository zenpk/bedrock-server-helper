import React from "react";
import ReactDOM from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { World } from "./page/World";
import { Home } from "./page/Home";
import "@mantine/core/styles.css";
import "./global.css";
import { createTheme, MantineProvider } from "@mantine/core";
import { authorization, redirectLogin } from "./util/myoauth.ts";
import { STORAGE_ACCESS_TOKEN } from "./util/constants.ts";

const theme = createTheme({});

const router = createBrowserRouter([
  {
    path: "/",
    element: <Home />,
  },
  {
    path: "/:worldId",
    element: <World />,
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
      <MantineProvider theme={theme}>
        <RouterProvider router={router} />
      </MantineProvider>
    </React.StrictMode>,
  );
}
