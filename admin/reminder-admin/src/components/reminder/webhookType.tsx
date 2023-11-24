import React from 'react'
import { Create, SimpleForm, TextInput, DateInput } from 'react-admin'
import { RichTextInput } from 'ra-input-rich-text';

export const WebhookTypes = [
    { id: "google_chat", name: "Google Chat" },
    { id: "discord", name: "Discord" },
    { id: "slack", name: "Slack" },
];
