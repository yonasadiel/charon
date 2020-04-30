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
  ROUTE_LOGIN,
  ROUTE_EVENT,
  ROUTE_EVENT_OVERVIEW,
  ROUTE_EVENT_PARTICIPATION,
  ROUTE_EVENT_QUESTION,
  ROUTE_EVENT_QUESTION_DETAIL,
  ROUTE_EVENT_QUESTION_EDIT,
  ROUTE_EVENT_QUESTION_EDIT_CREATE,
  ROUTE_EVENT_SYNC,
  ROUTE_EVENT_DECRYPT,
} from '../../routes';
import DecryptionPage from './decrypt/DecryptionPage';
import { hasPermissionForMenu, menuByRole } from './menu';
import ParticipationPage from './participation/ParticipationPage';
import QuestionPage from './question/QuestionPage';
import EventQuestionCreatePage from './question-editor/create/QuestionCreatePage';
import QuestionEditorPage from './question-editor/QuestionEditorPage';
import SynchronizationPage from './sync/SynchronizationPage';
import './EventPage.scss';

interface EventDetailPageProps extends RouteComponentProps<{ eventSlug: string }> {
  event: Event | null;
  participationKey: string | undefined;
  user: User | null;
  createQuestion: (eventSlug: string, question: Question) => Promise<void>;
  getQuestionsOfEvent: (eventSlug: string) => Promise<void>;
  verifyParticipation: (eventSlug: string, participationKey: string) => Promise<void>;
};


const renderEventDetailLoadingPage = (pathname: string, eventSlug: string) => (
  <div className="event-detail-page">
    <div className="titlebar">
      <h1><span className="skeleton">Judul ujian</span></h1>
    </div>
    <div className="menubar">
      <EventNavigation currentPath={pathname} eventSlug={eventSlug} menus={[]} />
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
    participationKey,
    user,
    verifyParticipation,
  } = props;

  React.useEffect(() => { document.title = event?.title || 'Ujian'; }, [event]);

  if (!event) return renderEventDetailLoadingPage(pathname, eventSlug);
  if (!user) return <Redirect to={ROUTE_LOGIN} />;

  const eventDetailLink = generateUrlWithParams(ROUTE_EVENT_OVERVIEW, { eventSlug: event.slug });
  const questionDetailLink = generateUrlWithParams(ROUTE_EVENT_QUESTION_DETAIL, { eventSlug: event.slug, questionNumber: 1 });
  const handleSubmitParticipationKey = (participationKey: string): Promise<void> => {
    return verifyParticipation(eventSlug, participationKey);
  };

  return (
    <div className="event-detail-page">
      <div className="titlebar">
        <h1>{event.title}</h1>
      </div>
      <div className="menubar">
        <EventNavigation currentPath={pathname} eventSlug={eventSlug} menus={menuByRole[user.role]} />
      </div>
      <div className="content">
        <Switch>
          <Route path={ROUTE_EVENT_QUESTION_EDIT_CREATE}>
            {hasPermissionForMenu(user, ROUTE_EVENT_QUESTION_EDIT_CREATE)
              ? <EventQuestionCreatePage event={event} />
              : <Redirect to={questionDetailLink} />}
          </Route>
          <Route path={ROUTE_EVENT_QUESTION_EDIT}>
            {hasPermissionForMenu(user, ROUTE_EVENT_QUESTION_EDIT)
              ? <QuestionEditorPage eventSlug={eventSlug} />
              : <Redirect to={questionDetailLink} />}
          </Route>
          <Route path={ROUTE_EVENT_QUESTION_DETAIL}>
            <QuestionPage user={user} />
          </Route>
          <Route path={ROUTE_EVENT_QUESTION}>
            <Redirect to={questionDetailLink} />
          </Route>
          <Route path={ROUTE_EVENT_PARTICIPATION}>
            {hasPermissionForMenu(user, ROUTE_EVENT_PARTICIPATION)
                ? <ParticipationPage />
                : <Redirect to={eventDetailLink} />}
          </Route>
          <Route path={ROUTE_EVENT_SYNC}>
            {hasPermissionForMenu(user, ROUTE_EVENT_SYNC)
                ? <SynchronizationPage />
                : <Redirect to={eventDetailLink} />}
          </Route>
          <Route path={ROUTE_EVENT_DECRYPT}>
            {hasPermissionForMenu(user, ROUTE_EVENT_DECRYPT)
                ? <DecryptionPage />
                : <Redirect to={eventDetailLink} />}
          </Route>
          <Route path={ROUTE_EVENT_OVERVIEW}>
            <EventDetail
              event={event}
              showParticipationKeyPrompt={(user.role === USER_ROLE.PARTICIPANT) && (participationKey === undefined)}
              onSubmitParticipatioNkey={handleSubmitParticipationKey}/>
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
  participationKey: charonSessionSelectors.getParticipationKey(state, props.match.params.eventSlug),
  user: charonSessionSelectors.getUser(state),
});

const mapDispatchToProps = {
  createQuestion: charonExamActions.createQuestion,
  getQuestionsOfEvent: charonExamActions.getQuestionsOfEvent,
  verifyParticipation: charonExamActions.verifyParticipation,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(EventDetailPage));
