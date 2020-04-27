import React from 'react';
import { Button, Card } from 'react-hephaestus';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';

import * as charonExamActions from '../../../../modules/charon/exam/action';
import * as charonExamSelectors from '../../../../modules/charon/exam/selector';
import { Question } from '../../../../modules/charon/exam/api';
import { generateUrlWithParams } from '../../../../modules/util/routes';
import { AppState } from '../../../../modules/store';
import { ROUTE_EVENT_QUESTION_EDIT_CREATE } from '../../../routes';
import QuestionEdit from './QuestionEdit';
import QuestionLoading from './QuestionEditLoading';
import './QuestionEditorPage.scss';

interface QuestionEditorPageProps {
  eventSlug: string;
}

interface ConnectedQuestionEditorPageProps extends QuestionEditorPageProps {
  deleteQuestion: (eventSlug: string, questionId: number) => Promise<void>;
  getQuestionsOfEvent: (eventSlug: string) => void;
  questions: Question[] | null;
}

const QuestionEditorPage = (props: ConnectedQuestionEditorPageProps) => {
  const { deleteQuestion, eventSlug, getQuestionsOfEvent, questions } = props;

  React.useEffect(() => { if (!questions) getQuestionsOfEvent(eventSlug); }, [getQuestionsOfEvent, eventSlug, questions]);
  const [isDeleting, setIsDeleting] = React.useState(false);
  const urlParam = { eventSlug };
  const handleDeleteQuestion = (questionId: number) => () => {
    setIsDeleting(true);
    deleteQuestion(eventSlug, questionId).then(() => {
      getQuestionsOfEvent(eventSlug);
      setIsDeleting(false);
    });
  };

  if (!questions || isDeleting) {
    return (
      <Card className="question-editor">
        {[0, 1, 2].map((idx) => <QuestionLoading key={idx} />)}
      </Card>
    );
  }
  return (
    <Card className="question-editor">
      <div className="add-row">
        <div>{questions.length} soal terdaftar.</div>
        <Link to={generateUrlWithParams(ROUTE_EVENT_QUESTION_EDIT_CREATE, urlParam)}>
          <Button><i className="fas fa-plus" />&nbsp;<strong>TAMBAH</strong></Button>
        </Link>
      </div>
      {questions?.map((question) => (
        <QuestionEdit question={question} onDeleteQuestion={handleDeleteQuestion(question.id)} />
      ))}
    </Card>
  );
};

const mapStateToProps = (state: AppState, props: QuestionEditorPageProps) => ({
  questions: charonExamSelectors.getQuestions(state, props.eventSlug),
});

const mapDispatchToProps = {
  getQuestionsOfEvent: charonExamActions.getQuestionsOfEvent,
  deleteQuestion: charonExamActions.deleteQuestion,
};

export default connect(mapStateToProps, mapDispatchToProps)(QuestionEditorPage);
