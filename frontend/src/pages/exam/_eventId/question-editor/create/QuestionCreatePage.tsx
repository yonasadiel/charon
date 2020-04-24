import React from 'react';
import { Card } from 'react-hephaestus';
import { connect } from 'react-redux';
import { useHistory, withRouter, RouteComponentProps } from 'react-router-dom';

import QuestionForm, { QuestionFormData } from '../../../../../components/exam/question/form/QuestionForm';
import { Event, Question } from '../../../../../modules/charon/exam/api';
import * as charonExamActions from '../../../../../modules/charon/exam/action';
import { generateUrlWithParams } from '../../../../../modules/util/routes';
import { ROUTE_EVENT_QUESTION_EDIT } from '../../../../routes';

interface ConnectedQuestionCreatePageProps extends RouteComponentProps<{ eventId: string }>, QuestionCreatePageProps {
  createQuestion: (eventId: number, question: Question) => Promise<void>;
  getQuestionsOfEvent: (eventId: number) => Promise<void>;
};

interface QuestionCreatePageProps {
  event: Event;
}

const QuestionCreatePage = (props: ConnectedQuestionCreatePageProps) => {
  const {
    event,
    getQuestionsOfEvent,
    createQuestion,
  } = props;

  React.useEffect(() => {
    document.title = !!event ? `Buat soal | ${event.title}` : 'Buat soal';
  }, [event]);
  const history = useHistory();

  if (!event) return <p>Loading</p>;

  const redirectAfterSubmitLink = generateUrlWithParams(ROUTE_EVENT_QUESTION_EDIT, { eventId: event.id });
  const handleSubmitNewQuestion = (questionData: QuestionFormData) => {
    createQuestion(event!.id, questionData as Question).then(() => {
      getQuestionsOfEvent(event.id);
      history.push(redirectAfterSubmitLink);
    });
  };

  return (
    <Card>
      <QuestionForm onSubmit={handleSubmitNewQuestion}/>
    </Card>
  );
};

const mapDispatchToProps = {
  createQuestion: charonExamActions.createQuestion,
  getQuestionsOfEvent: charonExamActions.getQuestionsOfEvent,
};

export default withRouter(connect(undefined, mapDispatchToProps)(QuestionCreatePage));
