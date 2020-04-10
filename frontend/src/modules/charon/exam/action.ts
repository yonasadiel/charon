import { AxiosError, AxiosResponse } from 'axios';

import { AppThunk } from '../../store';
import { CharonAPIError, CharonFormError } from '../http';
import { Event, Question } from './api';

export const PUT_EVENTS = 'charon/exam/PUT_EVENTS';
export const putEvents = (events: Event[] | null) => ({
  type: PUT_EVENTS,
  events,
});

export const PUT_QUESTIONS = 'charon/exam/PUT_QUESTIONS';
export const putQuestions = (eventId: number, questions: Question[] | null) => ({
  type: PUT_QUESTIONS,
  eventId,
  questions,
});

export function getEvents(): AppThunk<Promise<void>> {
  return async function (dispatch, _, { charonExamApi }) {
    dispatch(putEvents(null));
    return charonExamApi.getEvents()
      .then((res: AxiosResponse) => {
        const events: Event[] = res.data;
        events.forEach((event) => {
          event.questions = null;
          event.startsAt = new Date(event.startsAt);
          event.endsAt = new Date(event.endsAt);
        });
        dispatch(putEvents(events));
      })
      .catch((err: AxiosError) => {
        throw new CharonAPIError(err);
      });
  };
};

export function createEvent(event: Event): AppThunk<Promise<void>> {
  return async function (_dispatch, _getState, { charonExamApi }) {
    return charonExamApi.createEvent(event)
      .then(() => {
        return;
      })
      .catch((err: AxiosError) => {
        throw new CharonFormError(err);
      });
  };
};

export function getQuestionsOfEvent(eventId: number): AppThunk<Promise<void>> {
  return async function (dispatch, _, { charonExamApi }) {
    dispatch(putQuestions(eventId, null));
    return charonExamApi.getQuestionsOfEvent(eventId)
      .then((res: AxiosResponse) => {
        const questions: Question[] = res.data;
        dispatch(putQuestions(eventId, questions));
      })
      .catch((err: AxiosError) => {
        throw new CharonAPIError(err);
      });
  };
};

export function createQuestion(eventId: number, question: Question): AppThunk<Promise<void>> {
  return async function (_dispatch, _getState, { charonExamApi }) {
    return charonExamApi.createQuestion(eventId, question)
      .then(() => {
        return;
      })
      .catch((err: AxiosError) => {
        throw new CharonFormError(err);
      });
  };
};

export function deleteQuestion(eventId: number, questionId: number): AppThunk<Promise<void>> {
  return async function (_dispatch, _getState, { charonExamApi }) {
    return charonExamApi.deleteQuestion(eventId, questionId)
      .then(() => {
        return;
      })
      .catch((err: AxiosError) => {
        throw new CharonFormError(err);
      });
  };
};

