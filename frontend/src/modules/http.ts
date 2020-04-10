import axios, { AxiosError } from 'axios';

const http = axios.create({
  timeout: 5000,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true,
});

export abstract class HTTPError {
  public err: AxiosError;

  constructor(err: AxiosError) {
    this.err = err;
  }

  /**
   * Human readable message of the error
   */
  abstract getMessage(): string;

  /**
   * HTTP Status Code of the response
   * 0 stands for network error
   */
  public getStatusCode(): number {
    return this.err.response?.status || 0;
  }
}

export default http;
