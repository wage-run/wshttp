interface GoFetch {
  (req: Request, init?: RequestInit): Promise<Response>;
}
interface GoFetchConfig {
  max_retry?: number;
  // smux
  Version: 1 | 2;
  KeepAliveDisabled: boolean;
  KeepAliveInterval: string;
  KeepAliveTimeout: string;
  MaxFrameSize: number;
  MaxReceiveBuffer: number;
  MaxStreamBuffer: number;
}
declare global {
  var wshttpGen: (endpoint: string, config?: GoFetchConfig) => GoFetch;
}

export const GoFetchInit: Promise<void>;
