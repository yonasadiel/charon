import { AxiosError, AxiosResponse } from 'axios';
import { WordArray } from 'crypto-js';
import AES from 'crypto-js/aes';
import Base64 from 'crypto-js/enc-base64';
import Utf8 from 'crypto-js/enc-utf8';
import Hex from 'crypto-js/enc-hex';
import ModeCFB from 'crypto-js/mode-cfb';
import NoPadding from 'crypto-js/pad-nopadding';
import Sha256 from 'crypto-js/sha256';

import { AppThunk } from '../../store';
import { CharonAPIError, CharonFormError } from '../http';
import { Event, Participation, ParticipationStatus, Question, Venue, SynchronizationData } from './api';
import { getQuestions as getQuestionsSelector } from './selector';
import { putParticipationKey } from '../../session/action';

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

export const PUT_PARTICIPATION_STATUS = 'charon/exam/PUT_PARTICIPATION_STATUS';
export const putParticipationStatus = (eventSlug: string, participationStatus: ParticipationStatus[] | null) => ({
  type: PUT_PARTICIPATION_STATUS,
  eventSlug,
  participationStatus,
});

export const PUT_QUESTIONS = 'charon/exam/PUT_QUESTIONS';
export const putQuestions = (eventSlug: string, questions: Question[] | null) => ({
  type: PUT_QUESTIONS,
  eventSlug,
  questions,
});

export const PUT_QUESTION = 'charon/exam/PUT_QUESTION';
export const putQuestion = (eventSlug: string, question: Question) => ({
  type: PUT_QUESTION,
  eventSlug,
  question,
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

export function verifyParticipation(eventSlug: string, participationKey: string): AppThunk<Promise<void>> {
  return async function (dispatch, _, { charonExamApi }) {
    const hashedOnce: string = Sha256(participationKey).toString(Hex);
    return charonExamApi.verifyParticipation(eventSlug, hashedOnce)
      .then(() => {
        dispatch(putParticipationKey(eventSlug, participationKey));
      })
      .catch((err: AxiosError) => {
        throw new CharonFormError(err);
      });
  };
};

export function getParticipationStatus(eventSlug: string): AppThunk<Promise<void>> {
  return async function (dispatch, _, { charonExamApi }) {
    dispatch(putParticipationStatus(eventSlug, null));
    return charonExamApi.getParticipationStatus(eventSlug)
      .then((res: AxiosResponse) => {
        const participationStatus: ParticipationStatus[] = res.data;
        participationStatus.forEach((status) => {
          status.loginAt = new Date(status.loginAt);
        });
        dispatch(putParticipationStatus(eventSlug, participationStatus));
      })
      .catch((err: AxiosError) => {
        throw new CharonAPIError(err);
      });
  };
};

export function deleteParticipationStatus(eventSlug: string, sessionId: number): AppThunk<Promise<void>> {
  return async function (dispatch, _, { charonExamApi }) {
    dispatch(putParticipationStatus(eventSlug, null));
    return charonExamApi.deleteParticipationStatus(eventSlug, sessionId)
      .then(() => {
        return;
      })
      .catch((err: AxiosError) => {
        throw new CharonAPIError(err);
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

export function submitSubmission(eventSlug: string, participationKey: string, questionNumber: number, answer: string): AppThunk<Promise<void>> {
  return async function (dispatch, _, { charonExamApi }) {
    const encryptedPassword = encryptText(answer, participationKey);
    return charonExamApi.submitSubmission(eventSlug, questionNumber, encryptedPassword)
      .then((res: AxiosResponse<Question>) => {
        const question: Question = res.data;
        dispatch(putQuestion(eventSlug, question));
      })
      .catch((err: AxiosError) => {
        throw new CharonAPIError(err);
      });
  }
}

export function decryptEvent(eventSlug: string, key: string): AppThunk<Promise<void>> {
  return async function (_dispatch, _, { charonExamApi }) {
    return charonExamApi.decryptEvent(eventSlug, key)
      .then(() => { })
      .catch((err: AxiosError) => {
        throw new CharonAPIError(err);
      });
  };
};

function randomHex(length: number): string {
  let rnd = '';
  for (let i = 0; i < 2 * length; i++) {
    rnd += Math.floor((Math.random() * 16)).toString(16);
  }
  return rnd;
}

export function encryptText(plaintext: string, key: string): string {
  var keyBytes = Utf8.parse(key);
  var iv = Hex.parse(randomHex(16));
  var ciphertext = AES.encrypt(plaintext, keyBytes, {
      mode: ModeCFB,
      iv: iv,
      padding: NoPadding,
  });
  return ciphertext.iv.toString() + ciphertext.ciphertext.toString();
}

export function decryptText(ciphertext: string, key: string): string {
  const cipherHex = Base64.parse(ciphertext).toString(Hex);
  return decryptHex(cipherHex, key)
}

export function decryptHex(cipherHex: string, key: string): string {
  const keyBytes = Utf8.parse(key);
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
