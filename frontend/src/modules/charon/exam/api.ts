import { AxiosResponse } from 'axios';

import { User } from '../auth/api';
import conf from '../../../conf';
import http from '../../http';

export interface Event {
  id: number;
  slug: string;
  title: string;
  description: string;
  startsAt: Date;
  endsAt: Date;
  isDecrypted: boolean;
  lastSynchronization: Date;
  questions: { [number: number]: Question } | null;
  participations: { [id: number]: Participation } | null;
};

export interface Question {
  number: number;
  content: string;
  choices: string[];
  answer?: string;
};

export interface Participation {
  id: number;
  userUsername: string;
  venueId: number;
};

export interface Venue {
  id: number;
  name: string;
};

export interface SynchronizationData {
  event: Event;
  venue: Venue;
  questions: Question[];
  users: User[];
}

export interface CharonExamApi {
  getVenues: () => Promise<AxiosResponse<Venue[]>>;
  createVenue: (venue: Venue) => Promise<AxiosResponse<void>>;

  getEvents: () => Promise<AxiosResponse<Event[]>>;
  createEvent: (event: Event) => Promise<AxiosResponse<void>>;

  getParticipationsOfEvent: (eventSlug: string) => Promise<AxiosResponse<Participation[]>>;
  createParticipation: (eventSlug: string, participation: Participation) => Promise<AxiosResponse<Participation>>;
  verifyParticipation: (eventSlug: string, eventKey: string) => Promise<AxiosResponse<void>>;

  getQuestionsOfEvent: (eventSlug: string) => Promise<AxiosResponse<Question[]>>;
  createQuestion: (eventSlug: string, question: Question) => Promise<AxiosResponse<void>>;
  deleteQuestion: (eventSlug: string, questionNumber: number) => Promise<AxiosResponse<Question>>;

  getSynchronizationData: (eventSlug: string) => Promise<AxiosResponse<SynchronizationData>>;
  putSynchronizationData: (eventSlug: string, syncData: SynchronizationData) => Promise<AxiosResponse<SynchronizationData>>;

  decryptEvent: (eventSlug: string, key: string) => Promise<AxiosResponse<void>>;

  submitSubmission: (eventSlug: string, questionNumber: number, answer: string) => Promise<AxiosResponse<Question>>;
};

export default {
  getVenues: () => http.get(`${conf.charonApiUrl}/exam/venue/`),
  createVenue: (venue: Venue) => http.post(`${conf.charonApiUrl}/exam/venue/`, venue),

  getEvents: () => http.get(`${conf.charonApiUrl}/exam/`),
  createEvent: (event: Event) => http.post(`${conf.charonApiUrl}/exam/`, event),

  getParticipationsOfEvent: (eventSlug: string) => http.get(`${conf.charonApiUrl}/exam/${eventSlug}/participation/`),
  createParticipation: (eventSlug: string, participation: Participation) => http.post(`${conf.charonApiUrl}/exam/${eventSlug}/participation/`, participation),
  verifyParticipation: (eventSlug: string, eventKey: string) => http.post(`${conf.charonApiUrl}/exam/${eventSlug}/verify/`, { key: eventKey }),

  getQuestionsOfEvent: (eventSlug: string) => http.get(`${conf.charonApiUrl}/exam/${eventSlug}/question/`),
  createQuestion: (eventSlug: string, question: Question) => http.post(`${conf.charonApiUrl}/exam/${eventSlug}/question/`, question),
  deleteQuestion: (eventSlug: string, questionNumber: number) => http.delete(`${conf.charonApiUrl}/exam/${eventSlug}/question/${questionNumber}/`),

  getSynchronizationData: (eventSlug: string) => http.get(`${conf.charonApiUrl}/exam/${eventSlug}/sync/`),
  putSynchronizationData: (eventSlug: string, syncData: SynchronizationData) => http.post(`${conf.charonApiUrl}/exam/${eventSlug}/sync/`, syncData),

  decryptEvent: (eventSlug: string, key: string) => http.post(`${conf.charonApiUrl}/exam/${eventSlug}/decrypt/`, { key }),

  submitSubmission: (eventSlug: string, questionNumber: number, answer: string) => http.post(`${conf.charonApiUrl}/exam/${eventSlug}/question/${questionNumber}/submit/`, { answer }),
} as CharonExamApi;
