import React from "react";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { AlertCircle } from "lucide-react";

export function Card({
  h2 = "",
  children,
  alertText = "",
}: {
  h2?: string;
  children: React.ReactNode;
  alertText?: string;
}) {
  return (
    <div className={"card"}>
      {alertText && (
        <Alert>
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>Info</AlertTitle>
          <AlertDescription>{alertText}</AlertDescription>
        </Alert>
      )}
      <h1>Bedrock Server Helper</h1>
      <h2>{h2}</h2>
      {children}
    </div>
  );
}
