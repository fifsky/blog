export class AppError extends Error {
  code: string;
  details?: Record<string, string>;
  constructor(code: string, msg: string, details?: Record<string, string>) {
    super(msg);
    this.code = code;
    this.details = details;
    this.name = "AppError";
  }
}
