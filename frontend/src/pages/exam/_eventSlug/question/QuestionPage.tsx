import parseInt from 'lodash/parseInt';
import React from 'react';
import { Card } from 'react-hephaestus';
import { connect } from 'react-redux';
import { withRouter, RouteComponentProps } from 'react-router-dom';

import { USER_ROLE, User } from '../../../../modules/charon/auth/api';
import * as charonExamActions from '../../../../modules/charon/exam/action';
import * as charonExamSelectors from '../../../../modules/charon/exam/selector';
import { Event, Question } from '../../../../modules/charon/exam/api';
import * as charonSessionSelectors from '../../../../modules/session/selector';
import { AppState } from '../../../../modules/store';
import QuestionNavigation from '../../../../components/exam/question/navigation/QuestionNavigation';
import QuestionDetail from '../../../../components/exam/question/QuestionDetail';

interface QuestionPageProps extends RouteComponentProps<{ eventSlug: string, questionNumber: string }> {
  user: User;
}

interface ConnectedQuestionPageProps extends QuestionPageProps {
  event: Event | null;
  participationKey: string | undefined;
  questions: Question[] | null;
  currentQuestion: Question | null;
  getQuestionsOfEvent: (eventSlug: string) => Promise<void>;
  submitSubmission: (eventSlug: string, participationKey: string, questionNumber: number, answer: string) => Promise<void>;
}

const QuestionPage = (props: ConnectedQuestionPageProps) => {
  const {
    currentQuestion,
    event,
    getQuestionsOfEvent,
    match: { params: { questionNumber, eventSlug } },
    participationKey,
    questions,
    submitSubmission,
    user,
  } = props;

  const now = new Date();
  React.useEffect(() => { if (!questions) getQuestionsOfEvent(eventSlug); }, [getQuestionsOfEvent, eventSlug, questions]);

  if (!event) {
    return (
      <div className="question-pane">
        <p>Loading...</p>
      </div>
    );
  }

  if (user.role === USER_ROLE.PARTICIPANT && (event.startsAt > now || !event.isDecrypted)) {
    return (
      <Card className="question-pane">
        <p>Ujian belum dimulai. Soal akan muncul di sini jika ujian sudah dimulai.</p>
      </Card>
    );
  }

  if (user.role === USER_ROLE.PARTICIPANT && (!participationKey)) {
    return (
      <Card className="question-pane">
        <p>Masukkan password ujian dahulu sebelum memulai mengerjakan soal.</p>
      </Card>
    );
  }

  const handleSubmitSubmission = (answer: string) => {
    return submitSubmission(eventSlug, participationKey || '', parseInt(questionNumber), answer)
  }
  const initialAnswer = (user.role === USER_ROLE.PARTICIPANT && participationKey && !!currentQuestion && !!currentQuestion.answer)
    ? charonExamActions.decryptHex(currentQuestion.answer, participationKey)
    : '';

  return (
    <Card className="question-pane">
      <QuestionNavigation eventSlug={eventSlug} questions={questions} currentQuestionNumber={parseInt(questionNumber)} />
      <hr />
      <QuestionDetail
        initialAnswer={initialAnswer}
        question={currentQuestion}
        canAnswer={user.role === USER_ROLE.PARTICIPANT}
        onSubmitSubmission={handleSubmitSubmission} />
    </Card>
  );
};

const mapStateToProps = (state: AppState, props: RouteComponentProps<{ eventSlug: string, questionNumber: string }>) => ({
  event: charonExamSelectors.getEvent(state, props.match.params.eventSlug),
  participationKey: charonSessionSelectors.getParticipationKey(state, props.match.params.eventSlug),
  questions: charonExamSelectors.getQuestions(state, props.match.params.eventSlug),
  currentQuestion: charonExamSelectors.getQuestionByNumber(state, props.match.params.eventSlug, parseInt(props.match.params.questionNumber)),
});

const mapDispatchToProps = {
  getQuestionsOfEvent: charonExamActions.getQuestionsOfEvent,
  submitSubmission: charonExamActions.submitSubmission,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(QuestionPage));
