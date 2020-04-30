import React from 'react';
import { Button, TextInput } from 'react-hephaestus';

import { Question } from '../../../modules/charon/exam/api';
import './QuestionDetail.scss';

interface QuestionDetailProps {
  question: Question | null;
  canAnswer: boolean;
  initialAnswer: string;
  onSubmitSubmission: (answer: string) => Promise<void>;
};

const QuestionDetail = (props: QuestionDetailProps) => {
  const { canAnswer, question, onSubmitSubmission, initialAnswer } = props;

  const [answer, setAnswer] = React.useState(initialAnswer);
  const handleSubmitSubmission = () => {
    onSubmitSubmission(answer);
  };
  React.useEffect(() => {
    setAnswer(initialAnswer);
  }, [setAnswer, initialAnswer]);

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
      <div className="answer">
        {question.choices.length > 0
          ? question.choices.map((choice) => (
            <div className="choice" key={choice}>
              <input type="radio" value={choice} name="choice" onChange={() => setAnswer(choice)} checked={answer === choice} /> {choice}
            </div>
          ))
          : <TextInput placeholder="jawab..." value={answer} onChange={(e) => setAnswer(e.currentTarget.value)}/>}
      </div>
      {canAnswer && (
        <Button onClick={handleSubmitSubmission} className="submit">
          <strong>JAWAB</strong>
        </Button>
      )}
    </div>
  )
};

QuestionDetail.defaultProps = {
  canAnswer: false,
  onSubmitSubmission: () => Promise.reject,
  initialAnswer: '',
};

export default QuestionDetail;
