import { SubmissionError } from 'redux-form';

import { HTTPError } from '../http';

export class CharonAPIError extends HTTPError {

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
