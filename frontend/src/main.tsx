import React from "react";
import ReactDOM from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { World } from "./page/World";
import { Home } from "./page/Home";
import "@mantine/core/styles.css";
import "./global.css";
import { createTheme, MantineProvider } from "@mantine/core";

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

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <MantineProvider theme={theme}>
      <RouterProvider router={router} />
    </MantineProvider>
  </React.StrictMode>,
);
