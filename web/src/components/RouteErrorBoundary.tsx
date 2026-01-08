import { AlertCircle, Home } from "lucide-react";
import { isRouteErrorResponse, useRouteError, type ErrorResponse } from "react-router";
import { Button } from "./ui/button";

export function RouteErrorBoundary() {
  const error = useRouteError();

  let errorMessage = "An unexpected error occurred";
  let errorTitle = "Something went wrong";
  let errorDetails: string | null = null;

  if (isRouteErrorResponse(error)) {
    const response = error as ErrorResponse;
    switch (response.status) {
      case 404:
        errorTitle = "Page Not Found";
        errorMessage = "The page you're looking for doesn't exist or has been moved.";
        break;
      case 401:
        errorTitle = "Unauthorized";
        errorMessage = "You need to log in to access this page.";
        break;
      case 403:
        errorTitle = "Forbidden";
        errorMessage = "You don't have permission to access this page.";
        break;
      case 500:
        errorTitle = "Server Error";
        errorMessage = "Something went wrong on our end. Please try again later.";
        break;
      default:
        errorTitle = `Error ${response.status}`;
        errorMessage = response.statusText || errorMessage;
    }
    errorDetails = response.data?.toString() || null;
  } else if (error instanceof Error) {
    errorMessage = error.message;
    errorDetails = error.stack || null;
  }

  const handleGoHome = () => {
    window.location.href = "/";
  };

  const handleReload = () => {
    window.location.reload();
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-background">
      <div className="max-w-md w-full p-6 space-y-4">
        <div className="flex items-center gap-3 text-destructive">
          <AlertCircle className="w-8 h-8" />
          <h1 className="text-2xl font-bold">{errorTitle}</h1>
        </div>

        <p className="text-foreground/70">{errorMessage}</p>

        {errorDetails && (
          <details className="bg-muted p-3 rounded-md text-sm">
            <summary className="cursor-pointer font-medium">Error details</summary>
            <pre className="whitespace-pre-wrap break-words text-xs text-foreground/60">{errorDetails}</pre>
          </details>
        )}

        <div className="flex gap-2">
          <Button onClick={handleGoHome} variant="outline" className="flex-1 gap-2">
            <Home className="w-4 h-4" />
            Go Home
          </Button>
          <Button onClick={handleReload} className="flex-1 gap-2">
            Reload
          </Button>
        </div>
      </div>
    </div>
  );
}
