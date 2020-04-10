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

interface QuestionPaneProps extends RouteComponentProps<{ questionId: string }> {
  eventId: number;
}

interface ConnectedQuestionPaneProps extends QuestionPaneProps {
  event: Event | null;
  getQuestionsOfEvent: (eventId: number) => Promise<void>;
  questions: Question[] | null;
}

const QuestionPane = (props: ConnectedQuestionPaneProps) => {
  const { event, eventId, getQuestionsOfEvent, questions, match: { params: { questionId } } } = props;

  const now = new Date();
  React.useEffect(() => { if (!questions) getQuestionsOfEvent(eventId); }, [getQuestionsOfEvent, eventId, questions]);

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
      <QuestionNavigation eventId={eventId} questions={questions} currentQuestionId={parseInt(questionId)} />
      <hr />
      <QuestionDetail question={!!questions ? questions[0] : null} />
    </Card>
  );
};

const mapStateToProps = (state: AppState, props: RouteComponentProps<{ eventId: string, questionId: string }>) => ({
  event: charonExamSelectors.getEvent(state, parseInt(props.match.params.eventId)),
  questions: charonExamSelectors.getQuestions(state, parseInt(props.match.params.eventId)),
});

const mapDispatchToProps = {
  getQuestionsOfEvent: charonExamActions.getQuestionsOfEvent,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(QuestionPane));
