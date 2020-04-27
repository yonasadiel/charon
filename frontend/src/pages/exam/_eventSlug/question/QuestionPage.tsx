import parseInt from 'lodash/parseInt';
import React from 'react';
import { Card } from 'react-hephaestus';
import { connect } from 'react-redux';
import { withRouter, RouteComponentProps } from 'react-router-dom';

import * as charonExamActions from '../../../../modules/charon/exam/action';
import * as charonExamSelectors from '../../../../modules/charon/exam/selector';
import { Event, Question } from '../../../../modules/charon/exam/api';
import { AppState } from '../../../../modules/store';
import QuestionNavigation from '../../../../components/exam/question/navigation/QuestionNavigation';
import QuestionDetail from '../../../../components/exam/question/QuestionDetail';

interface QuestionPageProps extends RouteComponentProps<{ eventSlug: string, questionId: string }> {
}

interface ConnectedQuestionPageProps extends QuestionPageProps {
  event: Event | null;
  getQuestionsOfEvent: (eventSlug: string) => Promise<void>;
  questions: Question[] | null;
}

const QuestionPage = (props: ConnectedQuestionPageProps) => {
  const { event, getQuestionsOfEvent, questions, match: { params: { questionId, eventSlug } } } = props;

  const now = new Date();
  React.useEffect(() => { if (!questions) getQuestionsOfEvent(eventSlug); }, [getQuestionsOfEvent, eventSlug, questions]);

  if (!event) {
    return (
      <div className="question-pane">
        <p>Loading...</p>
      </div>
    );
  }

  if (event.startsAt > now) {
    return (
      <Card className="question-pane">
        <p>Ujian belum dimulai. Soal akan muncul di sini jika ujian sudah dimulai.</p>
      </Card>
    );
  }

  return (
    <Card className="question-pane">
      <QuestionNavigation eventSlug={eventSlug} questions={questions} currentQuestionId={parseInt(questionId)} />
      <hr />
      <QuestionDetail question={!!questions ? questions[0] : null} />
    </Card>
  );
};

const mapStateToProps = (state: AppState, props: RouteComponentProps<{ eventSlug: string, questionId: string }>) => ({
  event: charonExamSelectors.getEvent(state, props.match.params.eventSlug),
  questions: charonExamSelectors.getQuestions(state, props.match.params.eventSlug),
});

const mapDispatchToProps = {
  getQuestionsOfEvent: charonExamActions.getQuestionsOfEvent,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(QuestionPage));
