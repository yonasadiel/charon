import React from 'react';
import { connect } from 'react-redux';
import { withRouter, RouteComponentProps, Switch, Route, Redirect } from 'react-router-dom';

import { User, USER_ROLE } from '../../../modules/charon/auth/api';
import EventDetail from '../../../components/exam/event/EventDetail';
import EventNavigation from '../../../components/exam/event/navigation/EventNavigation';
import { Event, Question } from '../../../modules/charon/exam/api';
import * as charonExamActions from '../../../modules/charon/exam/action';
import * as charonExamSelectors from '../../../modules/charon/exam/selector';
import * as charonSessionSelectors from '../../../modules/session/selector';
import { generateUrlWithParams } from '../../../modules/util/routes';
import { AppState } from '../../../modules/store';
import {
  ROUTE_EVENT,
  ROUTE_EVENT_OVERVIEW,
  ROUTE_EVENT_PARTICIPATION,
  ROUTE_EVENT_QUESTION,
  ROUTE_EVENT_QUESTION_DETAIL,
  ROUTE_EVENT_QUESTION_EDIT,
  ROUTE_EVENT_QUESTION_EDIT_CREATE,
} from '../../routes';
import ParticipationPage from './participation/ParticipationPage';
import EventQuestionCreatePage from './question-editor/create/QuestionCreatePage';
import QuestionEditorPage from './question-editor/QuestionEditorPage';
import QuestionPage from './question/QuestionPage';
import './EventPage.scss';

interface EventDetailPageProps extends RouteComponentProps<{ eventSlug: string }> {
  event: Event | null;
  user: User | null;
  createQuestion: (eventSlug: string, question: Question) => Promise<void>;
  getQuestionsOfEvent: (eventSlug: string) => Promise<void>;
};

const canEditEvent = (user: User) => (user.role === USER_ROLE.ADMIN || user.role === USER_ROLE.ORGANIZER);

const renderEventDetailLoadingPage = (pathname: string, eventSlug: string) => (
  <div className="event-detail-page">
    <div className="titlebar">
      <h1><span className="skeleton">Judul ujian</span></h1>
    </div>
    <div className="menubar">
      <EventNavigation currentPath={pathname} eventSlug={eventSlug} hasEditPermission={false} />
    </div>
    <div className="content">
      <EventDetail />
    </div>
  </div>
);

const EventDetailPage = (props: EventDetailPageProps) => {
  const {
    location: { pathname },
    match: { params: { eventSlug } },
    event,
    user,
  } = props;

  React.useEffect(() => { document.title = event?.title || 'Ujian'; }, [event]);

  if (!event) return renderEventDetailLoadingPage(pathname, eventSlug);

  const hasEditPermission = !!user ? canEditEvent(user) : false;
  const eventDetailLink = generateUrlWithParams(ROUTE_EVENT_OVERVIEW, { eventSlug: event.slug });
  const questionDetailLink = generateUrlWithParams(ROUTE_EVENT_QUESTION_DETAIL, { eventSlug: event.slug, questionNumber: 1 });

  return (
    <div className="event-detail-page">
      <div className="titlebar">
        <h1>{event.title}</h1>
      </div>
      <div className="menubar">
        <EventNavigation currentPath={pathname} eventSlug={eventSlug} hasEditPermission={hasEditPermission} />
      </div>
      <div className="content">
        <Switch>
          <Route path={ROUTE_EVENT_QUESTION_EDIT_CREATE}>
            {hasEditPermission
              ? <EventQuestionCreatePage event={event} />
              : <Redirect to={questionDetailLink} />}
          </Route>
          <Route path={ROUTE_EVENT_QUESTION_EDIT}>
            {hasEditPermission
              ? <QuestionEditorPage eventSlug={eventSlug} />
              : <Redirect to={questionDetailLink} />}
          </Route>
          <Route path={ROUTE_EVENT_QUESTION_DETAIL}>
            <QuestionPage />
          </Route>
          <Route path={ROUTE_EVENT_QUESTION}>
            <Redirect to={questionDetailLink} />
          </Route>
          <Route path={ROUTE_EVENT_PARTICIPATION}>
            <ParticipationPage />
          </Route>
          <Route path={ROUTE_EVENT_OVERVIEW}>
            <EventDetail event={event} />
          </Route>
          <Route path={ROUTE_EVENT}>
            <Redirect to={eventDetailLink} />
          </Route>
        </Switch>
      </div>
    </div>
  );
};

const mapStateToProps = (state: AppState, props: RouteComponentProps<{ eventSlug: string }>) => ({
  event: charonExamSelectors.getEvent(state, props.match.params.eventSlug),
  user: charonSessionSelectors.getUser(state),
});

const mapDispatchToProps = {
  createQuestion: charonExamActions.createQuestion,
  getQuestionsOfEvent: charonExamActions.getQuestionsOfEvent,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(EventDetailPage));
