import flatMap from 'lodash/flatMap';

import { AppState } from '../../store';
import { getCharonState } from '../selector';
import { Event, Participation, ParticipationStatus, Question, Venue } from './api';
import { CharonExamState } from './reducer';

export function getCharonExamState(state: AppState): CharonExamState {
  return getCharonState(state).exam;
}

export function getVenues(state: AppState): Venue[] | null {
  const { venues } = getCharonExamState(state);
  return !!venues ? flatMap(venues) : venues;
}

export function getEvents(state: AppState): Event[] | null {
  const { events } = getCharonExamState(state);
  return !!events ? flatMap(events) : events;
}

export function getEvent(state: AppState, eventSlug: string): Event | null {
  const events = getCharonExamState(state).events;
  if (!events) return null;
  return events[eventSlug] || null;
}

export function getParticipations(state: AppState, eventSlug: string): Participation[] | null {
  const event = getEvent(state, eventSlug);
  if (!event) return null;
  if (!event.participations) return null;
  return flatMap(event.participations);
}

export function getParticipationStatus(state: AppState, eventSlug: string): ParticipationStatus[] | null {
  const event = getEvent(state, eventSlug);
  if (!event) return null;
  if (!event.participationStatus) return null;
  return flatMap(event.participationStatus);
}

export function getQuestions(state: AppState, eventSlug: string): Question[] | null {
  const event = getEvent(state, eventSlug);
  if (!event) return null;
  if (!event.questions) return null;
  return flatMap(event.questions);
}

export function getQuestionByNumber(state: AppState, eventSlug: string, questionNumber: number): Question | null {
  const event = getEvent(state, eventSlug);
  if (!event) return null;
  if (!event.questions) return null;
  return event.questions[questionNumber];
}
