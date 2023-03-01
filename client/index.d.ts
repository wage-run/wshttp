declare global {
  var GoFetch: (req: Request, init?: RequestInit) => Promise<Response>;
}

export const GoFetchInit: Promise<void>;
