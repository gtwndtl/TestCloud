import { CandidateInterface } from "./ICandidate";

export interface ElectionInterface {

    ID?: number;
  
    Title?: string;

    Description?: string;

    Start_time?: Date;

    End_time?: Date;

    Status?: string;

    Candidate_id?: CandidateInterface[];
  
  }