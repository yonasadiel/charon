import keyBy from 'lodash/keyBy';

import { PUT_EVENTS, PUT_PARTICIPATIONS, PUT_QUESTIONS, PUT_VENUES } from './action';
import { Event, Venue } from './api';

export interface CharonExamState {
  venues: { [id: number]: Venue } | null,
  events: { [slug: string]: Event } | null,
};

const initialState: CharonExamState = {
  venues: null,
  events: null,
};

export function charonExamReducer (state: CharonExamState = initialState, action: any) {
  switch (action.type) {
    case PUT_VENUES: {
      return {
        ...state,
        venues: action.venues === null ? null : keyBy(action.venues, 'id'),
      };
    }
    case PUT_EVENTS: {
      return {
        ...state,
        events: action.events === null ? null : keyBy(action.events, 'slug'),
      };
    }
    case PUT_PARTICIPATIONS: {
      const oldEvent = !!state.events ? state.events[action.eventSlug] : null;
      const newEvent = Object.assign({}, !!oldEvent ? oldEvent : {} as Event);
      newEvent.participations = action.participations === null ? null : keyBy(action.participations, 'id');
      return {
        ...state,
        events: {
          ...state.events,
          [action.eventSlug]: newEvent,
        },
      };
    }
    case PUT_QUESTIONS: {
      const oldEvent = !!state.events ? state.events[action.eventSlug] : null;
      const newEvent = Object.assign({}, !!oldEvent ? oldEvent : {} as Event);
      newEvent.questions = action.questions === null ? null : keyBy(action.questions, 'number');
      return {
        ...state,
        events: {
          ...state.events,
          [action.eventSlug]: newEvent,
        },
      };
    }
    default:
      return state;
  }
}
