// filepath: c:\Users\lastb\Code\send_blue\frontend\send-blue-ui\src\types.d.ts

export interface Client {
    uid: string;
    name: string;
    messagesSent: number;
}

export interface ClientLead {
    leadNumber: string;
    clientUid: string;
    messagesReceived: number;
    lastContacted: Date | null;
}

export interface MessageQueue {
    msgUid: string;
    messageBody: string;
    fromClientId: string;
    toClientLead: string;
    scheduledSendTime: Date | null;
    timeSent: Date | null;
    status: string;
}

export interface MessageEventHistory {
    msgUid: string;
    time_stamp: Date;
    prev_status: string | null;
    curr_status: string;
}

export interface MessageArchive {
    msgUid: string;
    messageBody: string;
    from_client_id: string;
    to_client_lead: string;
    time_archived: Date | null;
    status: string;
}