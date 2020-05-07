import { SubmissionError } from 'redux-form';
import axios, { AxiosError } from 'axios';

import conf from '../../conf';

const charonApiHttp = axios.create({
  baseURL: conf.charonApiUrl,
  timeout: 5000,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true,
});

export default charonApiHttp;

export class CharonAPIError {
  public err: AxiosError;

  constructor(err: AxiosError) {
    this.err = err;
  }

  /**
   * Human readable message of the error
   */
  public getMessage(): string {
    return this.err.response?.data.message || "Unknown Error";
  }

  /**
   * Charon constant error code of the response
   */
  public getErrorCode(): string {
    return this.err.response?.data.code || "unknown_error";
  }
}

export class CharonFormError extends CharonAPIError {
  public getMessage(): string {
    return 'Error submitting response';
  }

  public asSubmissionError(): SubmissionError {
    return new SubmissionError(this.err.response?.data.message);
  }
}
