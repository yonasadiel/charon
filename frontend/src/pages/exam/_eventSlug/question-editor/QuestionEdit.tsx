import React from 'react';
import { Button } from 'react-hephaestus';

import { Question } from '../../../../modules/charon/exam/api';
import './QuestionEditorPage.scss';

interface QuestionEditProps {
  question: Question;
  onDeleteQuestion: () => void;
}

const QuestionEdit = (props: QuestionEditProps) => {
  const { onDeleteQuestion, question } = props;

  return (
    <div className="question">
      <div className="content" dangerouslySetInnerHTML={{ __html: question.content }}></div>
      <ul className="choices">
        {question.choices.map((choice) => (
          <li className="choice" key={choice}>{choice}</li>
        ))}
      </ul>
      <div className="editing-row">
        <Button buttonType="outlined" className="delete-button" onClick={onDeleteQuestion}>
          <i className="fas fa-trash"></i><strong>HAPUS</strong>
        </Button>
      </div>
    </div>
  );
};

export default QuestionEdit;