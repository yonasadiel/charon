import React from 'react';
import { Field, reduxForm, InjectedFormProps, WrappedFieldProps } from 'redux-form';
import { Button, TextInput, BUTTON_TYPE_CONTAINED } from 'react-hephaestus';

import { Venue } from '../../../../modules/charon/exam/api';
import LoadingCircle from '../../../loading/Circle';
import './VenueForm.scss';

export interface VenueFormData extends Partial<Venue> {};

const renderContentField = (field: WrappedFieldProps) => {
  const { input } = field;
  return (
    <TextInput placeholder="Lokasi baru" {...input} wrapperStyle={{flex: 1}} />
  );
};

const VenueForm = (props: InjectedFormProps<VenueFormData>) => {
  const { error, handleSubmit, submitting, reset } = props;
  return (
    <form className="venue-form" onSubmit={(e) => { handleSubmit(e); reset() }}>
      <div className="field-row">
        <Field name="name" component={renderContentField} />
        <Button buttonType={BUTTON_TYPE_CONTAINED} type="submit" className="button">
          { !submitting ? (<strong>BUAT</strong>) : <LoadingCircle /> }
        </Button>
      </div>
      <small className="error">{ error }</small>
    </form>
  );
};

export default reduxForm<VenueFormData>({ form: 'venue' })(VenueForm);
