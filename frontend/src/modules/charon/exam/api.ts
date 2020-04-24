import { AxiosResponse } from 'axios';

import conf from '../../../conf';
import http from '../../http';

export interface Event {
  id: number;
  title: string;
  description: string;
  startsAt: Date;
  endsAt: Date;
  questions: { [id: number]: Question } | null;
};

export interface Question {
  id: number;
  content: string;
  choices: string[];
  answer?: string;
};

export interface Venue {
  id: number;
  name: string;
}

export interface CharonExamApi {
  getVenues: () => Promise<AxiosResponse<Venue[]>>;
  createVenue: (venue: Venue) => Promise<AxiosResponse<void>>;

  getEvents: () => Promise<AxiosResponse<Event[]>>;
  createEvent: (event: Event) => Promise<AxiosResponse<void>>;

  getQuestionsOfEvent: (eventId: number) => Promise<AxiosResponse<Question[]>>;
  createQuestion: (eventId: number, question: Question) => Promise<AxiosResponse<void>>;
  deleteQuestion: (eventId: number, questionId: number) => Promise<AxiosResponse<Question>>;
};

export default {
  getVenues: () => http.get(`${conf.charonApiUrl}/exam/venue/`),
  createVenue: (venue: Venue) => http.post(`${conf.charonApiUrl}/exam/venue/`, venue),

  getEvents: () => http.get(`${conf.charonApiUrl}/exam/`),
  createEvent: (event: Event) => http.post(`${conf.charonApiUrl}/exam/`, event),

  getQuestionsOfEvent: (eventId: number) => http.get(`${conf.charonApiUrl}/exam/${eventId}/question/`),
  createQuestion: (eventId: number, question: Question) => http.post(`${conf.charonApiUrl}/exam/${eventId}/question/`, question),
  deleteQuestion: (eventId: number, questionId: number) => http.delete(`${conf.charonApiUrl}/exam/${eventId}/question/${questionId}/`),
} as CharonExamApi;
