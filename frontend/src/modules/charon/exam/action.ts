import { AxiosError, AxiosResponse } from 'axios';
import { WordArray } from 'crypto-js';
import AES from 'crypto-js/aes';
import Base64 from 'crypto-js/enc-base64';
import Utf8 from 'crypto-js/enc-utf8';
import Hex from 'crypto-js/enc-hex';
import ModeCFB from 'crypto-js/mode-cfb';
import NoPadding from 'crypto-js/pad-nopadding';

import { AppThunk } from '../../store';
import { CharonAPIError, CharonFormError } from '../http';
import { Event, Participation, Question, Venue, SynchronizationData } from './api';
import { getQuestions as getQuestionsSelector } from './selector';

export const PUT_VENUES = 'charon/exam/PUT_VENUES';
export const putVenues = (venues: Venue[] | null) => ({
  type: PUT_VENUES,
  venues,
});

export const PUT_EVENTS = 'charon/exam/PUT_EVENTS';
export const putEvents = (events: Event[] | null) => ({
  type: PUT_EVENTS,
  events,
});

export const PUT_PARTICIPATIONS = 'charon/exam/PUT_PARTICIPATIONS';
export const putParticipations = (eventSlug: string, participations: Participation[] | null) => ({
  type: PUT_PARTICIPATIONS,
  eventSlug,
  participations,
});

export const PUT_QUESTIONS = 'charon/exam/PUT_QUESTIONS';
export const putQuestions = (eventSlug: string, questions: Question[] | null) => ({
  type: PUT_QUESTIONS,
  eventSlug,
  questions,
});

export function getVenues(): AppThunk<Promise<void>> {
  return async function (dispatch, _, { charonExamApi }) {
    dispatch(putVenues(null));
    return charonExamApi.getVenues()
      .then((res: AxiosResponse) => {
        const venues: Venue[] = res.data;
        dispatch(putVenues(venues));
      })
      .catch((err: AxiosError) => {
        throw new CharonAPIError(err);
      });
  };
};

export function createVenue(venue: Venue): AppThunk<Promise<void>> {
  return async function (_dispatch, _getState, { charonExamApi }) {
    return charonExamApi.createVenue(venue)
      .then(() => {
        return;
      })
      .catch((err: AxiosError) => {
        throw new CharonFormError(err);
      });
  };
};

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

export function getParticipationsOfEvent(eventSlug: string): AppThunk<Promise<void>> {
  return async function (dispatch, _, { charonExamApi }) {
    dispatch(putParticipations(eventSlug, null));
    return charonExamApi.getParticipationsOfEvent(eventSlug)
      .then((res: AxiosResponse) => {
        const participations: Participation[] = res.data;
        dispatch(putParticipations(eventSlug, participations));
      })
      .catch((err: AxiosError) => {
        throw new CharonAPIError(err);
      });
  };
};

export function createParticipation(eventSlug: string, participation: Participation): AppThunk<Promise<void>> {
  return async function (_dispatch, _getState, { charonExamApi }) {
    return charonExamApi.createParticipation(eventSlug, participation)
      .then(() => {
        return;
      })
      .catch((err: AxiosError) => {
        throw new CharonFormError(err);
      });
  };
};

export function getQuestionsOfEvent(eventSlug: string): AppThunk<Promise<void>> {
  return async function (dispatch, _, { charonExamApi }) {
    dispatch(putQuestions(eventSlug, null));
    return charonExamApi.getQuestionsOfEvent(eventSlug)
      .then((res: AxiosResponse) => {
        const questions: Question[] = res.data;
        dispatch(putQuestions(eventSlug, questions));
      })
      .catch((err: AxiosError) => {
        throw new CharonAPIError(err);
      });
  };
};

export function createQuestion(eventSlug: string, question: Question): AppThunk<Promise<void>> {
  return async function (_dispatch, _getState, { charonExamApi }) {
    return charonExamApi.createQuestion(eventSlug, question)
      .then(() => {
        return;
      })
      .catch((err: AxiosError) => {
        throw new CharonFormError(err);
      });
  };
};

export function deleteQuestion(eventSlug: string, questionNumber: number): AppThunk<Promise<void>> {
  return async function (_dispatch, _getState, { charonExamApi }) {
    return charonExamApi.deleteQuestion(eventSlug, questionNumber)
      .then(() => {
        return;
      })
      .catch((err: AxiosError) => {
        throw new CharonFormError(err);
      });
  };
};

export function getSynchronizationData(eventSlug: string): AppThunk<Promise<SynchronizationData>> {
  return async function (_dispatch, _, { charonExamApi }) {
    return charonExamApi.getSynchronizationData(eventSlug)
      .then((res: AxiosResponse) => {
        const syncData: SynchronizationData = res.data;
        return syncData;
      })
      .catch((err: AxiosError) => {
        throw new CharonAPIError(err);
      });
  };
};

export function putSynchronizationData(eventSlug: string, syncData: SynchronizationData): AppThunk<Promise<void>> {
  return async function (dispatch, _, { charonExamApi }) {
    return charonExamApi.putSynchronizationData(eventSlug, syncData)
      .then(() => {
        dispatch(putEvents(null));
        dispatch(putQuestions(eventSlug, null));
        dispatch(putParticipations(eventSlug, null));
      })
      .catch((err: AxiosError) => {
        throw new CharonAPIError(err);
      });
  };
};

export function decryptEvent(eventSlug: string, key: string): AppThunk<Promise<void>> {
  return async function (_dispatch, _, { charonExamApi }) {
    return charonExamApi.decryptEvent(eventSlug, key)
      .then(() => { })
      .catch((err: AxiosError) => {
        throw new CharonAPIError(err);
      });
  };
};

function decryptText(ciphertext: string, key: string): string {
  const keyBytes = Utf8.parse(key);
  const cipherHex = Base64.parse(ciphertext).toString(Hex);
  const iv = Hex.parse(cipherHex.slice(0, 32));
  const cipherBytes = Hex.parse(cipherHex.slice(32));
  const plainBytes  = AES.decrypt({
    ciphertext: cipherBytes,
    salt: '',
  } as WordArray, keyBytes, {
    iv: iv,
    mode: ModeCFB,
    padding: NoPadding,
  });
  return plainBytes.toString(Utf8);
}

export function decryptEventLocal(eventSlug: string, key: string): AppThunk<Promise<void>> {
  return async function (dispatch, getState) {
    const questions = getQuestionsSelector(getState(), eventSlug);
    if (!!questions) {
      for (let i = 0; i < questions.length; i++) {
        questions[i].content = decryptText(questions[i].content, key);
        const choiceText = decryptText(questions[i].choices[0], key);
        const choiceTexts = choiceText.split('|');
        questions[i].choices = [];
        for (let j = 0; j < choiceTexts.length; j++) {
          if (choiceTexts[j].length > 0) {
            questions[i].choices.push(choiceTexts[j]);
          }
        }
      }
      dispatch(putQuestions(eventSlug, questions));
    }
  };
};
