import React from 'react';
import { Link } from 'react-router-dom';

import { Question } from '../../../../modules/charon/exam/api';
import { generateUrlWithParams } from '../../../../modules/util/routes';
import './QuestionNavigation.scss';
import { ROUTE_EVENT_QUESTION_DETAIL } from '../../../../pages/routes';

interface QuestionNavigationProps {
  currentQuestionNumber: number;
  eventSlug: string;
  questions: Question[] | null;
}

const QuestionNavigation = (props: QuestionNavigationProps) => {
  const { currentQuestionNumber, eventSlug, questions } = props;
  if (!questions) {
    return <p>Loading</p>;
  }
  return (
    <div className="question-navigation">
    {questions?.map((question, i) => (
      <Link
        className={`nav-item ${question.number === currentQuestionNumber ? 'active' : ''}`}
        to={generateUrlWithParams(ROUTE_EVENT_QUESTION_DETAIL, { eventSlug, questionNumber: i + 1 })}
        key={question.number}>
        <span>{i + 1}</span>
      </Link>
    ))}
    </div>
  );
};

export default QuestionNavigation;
