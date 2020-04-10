import flatMap from 'lodash/flatMap';

import { AppState } from '../../store';
import { getCharonState } from '../selector';
import { Event, Question } from './api';
import { CharonExamState } from './reducer';

export function getCharonExamState(state: AppState): CharonExamState {
  return getCharonState(state).exam;
}

export function getEvents(state: AppState): Event[] | null {
  const { events } = getCharonExamState(state);
  return !!events ? flatMap(events) : events;
}

export function getEvent(state: AppState, eventId: number): Event | null {
  const events = getCharonExamState(state).events;
  if (!events) return null;
  return events[eventId] || null;
}

export function getQuestions(state: AppState, eventId: number): Question[] | null {
  const event = getEvent(state, eventId);
  if (!event) return null;
  if (!event.questions) return null;
  return flatMap(event.questions);
}

export function getQuestionByNumber(state: AppState, eventId: number, questionNumber: number): Question | null {
  const questions = getQuestions(state, eventId);
  return !!questions ? questions[questionNumber - 1] : null;
}
