import React from 'react';
import { TextInput } from 'react-hephaestus';

import { Question } from '../../../modules/charon/exam/api';
import './QuestionDetail.scss';

interface QuestionDetailProps {
  question: Question | null;
};

const QuestionDetail = (props: QuestionDetailProps) => {
  const { question } = props;

  if (!question) {
    return (
      <div>
        <p><span className="skeleton">Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur justo massa, imperdiet et nulla id, lobortis egestas nulla. Cras non quam dui. Vestibulum rhoncus, ante eget faucibus suscipit, ante lacus congue nulla, id pretium velit purus quis sapien.</span></p>
        <ul>
          <li className="skeleton">Pilihan 1</li>
          <li className="skeleton">Pilihan 2</li>
          <li className="skeleton">Pilihan 3</li>
          <li className="skeleton">Pilihan 4</li>
        </ul>
      </div>
    )
  }
  return (
    <div className="question-detail">
      <p dangerouslySetInnerHTML={{ __html: question.content }}></p>
      <div className="choices">
        {question.choices.length > 0
          ? question.choices.map((choice) => (
            <div className="choice" key={choice}>
              <input type="radio" value={choice} name="choice" /> {choice}
            </div>
          ))
          : <TextInput placeholder="jawab..." />}
      </div>
    </div>
  )
};

export default QuestionDetail;
