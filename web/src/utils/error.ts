export class AppError extends Error {
  code: number;
  constructor(code: number, msg: string) {
    super(msg);
    this.code = code;
    this.name = "AppError";
  }
}
