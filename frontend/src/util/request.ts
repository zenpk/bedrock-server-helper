import axios from "axios";
import {
  BASE_URL,
  STORAGE_ACCESS_TOKEN,
  STORAGE_REFRESH_TOKEN,
} from "./constants.ts";
import { redirectLogin, refresh } from "./myoauth.ts";
import React, { SetStateAction } from "react";

export async function get(url: string) {
  if (!window.localStorage.getItem(STORAGE_ACCESS_TOKEN)) {
    redirectLogin();
  }
  try {
    const resp = await axios.get(BASE_URL + url, {
      headers: {
        Authorization: `Bearer ${window.localStorage.getItem(
          STORAGE_ACCESS_TOKEN,
        )}`,
      },
    });
    return resp.data;
  } catch (e: any) {
    handleError(e);
  }
}

export async function post(url: string, body: any) {
  if (!window.localStorage.getItem(STORAGE_ACCESS_TOKEN)) {
    redirectLogin();
  }
  try {
    const resp = await axios.post(BASE_URL + url, body, {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${window.localStorage.getItem(
          STORAGE_ACCESS_TOKEN,
        )}`,
      },
    });
    return resp.data;
  } catch (e: any) {
    handleError(e);
  }
}

export async function del(url: string, body: any) {
  if (!window.localStorage.getItem(STORAGE_ACCESS_TOKEN)) {
    redirectLogin();
  }
  try {
    const resp = await axios.delete(BASE_URL + url, {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${window.localStorage.getItem(
          STORAGE_ACCESS_TOKEN,
        )}`,
      },
      data: body,
    });
    return resp.data;
  } catch (e: any) {
    handleError(e);
  }
}

export async function getEvent(
  url: string,
  setResult: React.Dispatch<SetStateAction<string>>,
) {
  try {
    const response = await fetch(BASE_URL + url, {
      method: "GET",
      headers: {
        "Content-Type": "text/event-stream",
      },
    });
    return handleEventResp(response, setResult);
  } catch (e) {
    console.log(e);
    return null;
  }
}

export async function postEvent(
  url: string,
  body: any,
  setResult: React.Dispatch<SetStateAction<string>>,
) {
  try {
    const response = await fetch(BASE_URL + url, {
      method: "POST",
      headers: {
        "Content-Type": "text/event-stream",
      },
      body: body,
    });
    return handleEventResp(response, setResult);
  } catch (e) {
    console.log(e);
    return null;
  }
}

function handleError(e: any) {
  switch (e.response?.status) {
    case 401:
      break;
    default:
      console.log(e);
      break;
  }
  return null;
}

function handleUnauthorized() {
  if (window.localStorage.getItem(STORAGE_REFRESH_TOKEN)) {
    refresh().catch(() => {
      redirectLogin();
    });
  } else {
    redirectLogin();
  }
}

async function handleEventResp(
  resp: Response,
  setResult: React.Dispatch<SetStateAction<string>>,
) {
  const reader = resp.body?.pipeThrough(new TextDecoderStream()).getReader();
  if (resp.status === 401) {
    handleUnauthorized();
    return null;
  }
  while (reader) {
    const { value, done } = await reader.read();
    if (done) break;
    setResult((prev) => prev + value);
  }
}
