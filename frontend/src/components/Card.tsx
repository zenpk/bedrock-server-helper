import React from "react";

export function Card({
  h2 = "",
  children,
}: {
  h2?: string;
  children: React.ReactNode;
}) {
  return (
    <div className={"card"}>
      <h1>Bedrock Server Helper</h1>
      <h2>{h2}</h2>
      {children}
    </div>
  );
}
