import React from 'react'
import { Edit, SimpleForm, TextInput, DateInput } from 'react-admin'
import { RichTextInput } from 'ra-input-rich-text';

const TaskEdit = (props) => {
    return (
        <Edit title='Edit Task' {...props}>
            <SimpleForm>
                <TextInput disabled source='id' />
                <TextInput source='name' />
                <TextInput multiline source='schedule' />
                <TextInput multiline source='webhook' />
                <RichTextInput multiline source='message' />
            </SimpleForm>
        </Edit>
    )
}

export default TaskEdit