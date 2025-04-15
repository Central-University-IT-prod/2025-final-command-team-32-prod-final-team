import { Terminal } from "lucide-react";

import { Alert, AlertDescription, AlertTitle } from "@/components/ui/lib/alert";

export function CustomAlert({ title, message }) {
  return (
    <Alert className="mt-5">
      <Terminal className="h-4 w-4" />
      <AlertTitle className="text-xl">{title}</AlertTitle>
      <AlertDescription className="text-lg">{message}</AlertDescription>
    </Alert>
  );
}
