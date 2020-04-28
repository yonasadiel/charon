export const ROUTE_HOME = '/';

export const ROUTE_LOGIN = ROUTE_HOME + 'login/';

export const ROUTE_USER = ROUTE_HOME + 'user/'
export const ROUTE_USER_LIST = ROUTE_USER;

export const ROUTE_VENUE = ROUTE_HOME + 'venue/';
export const ROUTE_VENUE_LIST = ROUTE_VENUE;

export const ROUTE_EXAM = ROUTE_HOME + 'exam/';
export const ROUTE_EVENT_LIST = ROUTE_EXAM;
export const ROUTE_EVENT = ROUTE_EXAM + ':eventSlug/';
export const ROUTE_EVENT_OVERVIEW = ROUTE_EVENT + 'overview/';
export const ROUTE_EVENT_PARTICIPATION = ROUTE_EVENT + 'participation/';
export const ROUTE_EVENT_QUESTION = ROUTE_EVENT + 'question/';
export const ROUTE_EVENT_QUESTION_DETAIL = ROUTE_EVENT_QUESTION + ':questionNumber/';
export const ROUTE_EVENT_QUESTION_EDIT = ROUTE_EVENT + 'question-editor/';
export const ROUTE_EVENT_QUESTION_EDIT_CREATE = ROUTE_EVENT_QUESTION_EDIT + 'new/';
export const ROUTE_EVENT_SYNC = ROUTE_EVENT + 'sync/';
export const ROUTE_EVENT_DECRYPT = ROUTE_EVENT + 'decrypt/';
