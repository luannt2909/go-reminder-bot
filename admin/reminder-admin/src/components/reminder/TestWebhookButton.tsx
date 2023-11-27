import React, {useState, Fragment} from 'react'
import IconContentSend from "@material-ui/icons/Send";
import IconCancel from '@material-ui/icons/Cancel';
import {Dialog, DialogActions, DialogContent, DialogContentText, Toolbar} from "@material-ui/core";
import DialogTitle from "@material-ui/core/DialogTitle";
import {useCreate, TextInput, SaveButton, Form, Button, SimpleForm,
    useRefresh, SelectInput, fetchUtils} from "react-admin";
import {required, useNotify, email } from "ra-core";
import {WebhookTypes} from "./webhookType";

const apiURL = import.meta.env.VITE_SIMPLE_REST_URL

const WebhookTestButton = () => {
    const [showDialog, setShowDialog] = useState(false);
    const notify = useNotify();
    const refresh = useRefresh();

    const handleClick = () => {
        setShowDialog(true);
    };

    const handleCloseClick = () => {
        setShowDialog(false);
    };

    const handleSubmit = async values => {
        fetchUtils.fetchJson(`${apiURL}/webhook/send`, { method: 'POST', body: JSON.stringify(values) })
            .then(() => {
                notify('Send message successful');
            })
            .catch((e) => {
                notify(`Error: Send message failed: ${e.message}`, { type: 'error' })
            })
            .finally(() => {
                // setLoading(false);
            });
    };

    return (
        <>
            <Button onClick={handleClick} label="Test Webhook">
                <IconContentSend />
            </Button>
            <Dialog
                fullWidth
                open={showDialog}
                onClose={handleCloseClick}
                aria-label="Test Webhook"
            >
                <DialogTitle>Send Message Through Webhook</DialogTitle>
                <DialogContent>
                    <SimpleForm
                        resource="webhook"
                        // We override the redux-form onSubmit prop to handle the submission ourselves
                        onSubmit={handleSubmit}
                        // We want no toolbar at all as we have our modal actions
                        toolbar={<TestWebhookButtonToolbar onCancel={handleCloseClick}/>}
                            >
                        <SelectInput choices={WebhookTypes} required source='webhook_type'/>
                        <TextInput type="url" multiline fullWidth required source='webhook'/>
                        <TextInput multiline fullWidth required source='message'/>
                    </SimpleForm>

                </DialogContent>
            </Dialog>
        </>
    );
}

function TestWebhookButtonToolbar({ onCancel, ...props }) {
    return (
        <Toolbar {...props}>
            <SaveButton icon={<IconContentSend/>} label="Send" submitOnEnter={true} />
            <CloseButton onClick={onCancel} />
        </Toolbar>
    );
}

function CloseButton(props) {
    return (
        <Button label="Close" {...props}>
            <IconCancel />
        </Button>
    );
}

export default WebhookTestButton