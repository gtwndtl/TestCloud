import { UsersInterface } from "../../interfaces/IUser";

import { SignInInterface } from "../../interfaces/SignIn";

import axios from "axios";

const apiUrl = "http://localhost:8000";

const token = localStorage.getItem("token");
const tokenType = localStorage.getItem("token_type");

const requestOptions = {
  headers: {
    "Content-Type": "application/json",
    Authorization: `${tokenType} ${token}`,
  },
};


async function SignIn(data: SignInInterface) {

  return await axios

    .post("http://localhost/api/user/signin", data, requestOptions)

    .then((res) => res)

    .catch((e) => e.response);

}


async function GetGender() {

  return await axios

    .get(`${apiUrl}/genders`, requestOptions)

    .then((res) => res)

    .catch((e) => e.response);

}


async function GetUsers() {

  return await axios

    .get("http://localhost/api/user/users", requestOptions)

    .then((res) => res)

    .catch((e) => e.response);

}


async function GetUsersById(id: string) {

  return await axios

    .get(`http://localhost/api/user/user/${id}`, requestOptions)

    .then((res) => res)

    .catch((e) => e.response);

}


async function UpdateUsersById(id: string, data: UsersInterface) {

  return await axios

    .put(`http://localhost/api/user/user/${id}`, data, requestOptions)

    .then((res) => res)

    .catch((e) => e.response);

}


async function DeleteUsersById(id: string) {

  return await axios

    .delete(`http://localhost/api/user/user/${id}`, requestOptions)

    .then((res) => res)

    .catch((e) => e.response);

}


async function CreateUser(data: UsersInterface) {

  return await axios

    .post("http://localhost/api/user/signup", data, requestOptions)

    .then((res) => res)

    .catch((e) => e.response);

}

async function GetCandidates() {
  return await axios
    .get("http://localhost/api/candidate/candidates")
    .then((res) => res)
    .catch((e) => e.response);
}



async function GetCandidateById(id: string) {

  return await axios

    .get(`http://localhost/api/candidate/candidate/${id}`)

    .then((res) => res)

    .catch((e) => e.response);

}


async function UpdateCandidateById(id: string, data: UsersInterface) {

  return await axios

    .put(`http://localhost/api/candidate/candidate/${id}`, data)

    .then((res) => res)

    .catch((e) => e.response);

}


async function DeleteCandidateById(id: string) {

  return await axios

    .delete(`http://localhost/api/candidate/candidate/${id}`)

    .then((res) => res)

    .catch((e) => e.response);

}


async function CreateCandidate(data: UsersInterface) {

  return await axios

    .post("http://localhost/api/candidate/candidate", data)

    .then((res) => res)

    .catch((e) => e.response);

}



async function GetElections() {
  return await axios
    .get("http://localhost/api/election/elections")
    .then((res) => res)
    .catch((e) => e.response);
}



async function GetElectionById(id: string) {

  return await axios

    .get(`http://localhost/api/election/election/${id}`)

    .then((res) => res)

    .catch((e) => e.response);

}


async function UpdateElectionById(id: string, data: UsersInterface) {

  return await axios

    .put(`http://localhost/api/election/election/${id}`, data)

    .then((res) => res)

    .catch((e) => e.response);

}


async function DeleteElectionById(id: string) {

  return await axios

    .delete(`http://localhost/api/election/election/${id}`)

    .then((res) => res)

    .catch((e) => e.response);

}


async function CreateElection(data: UsersInterface) {

  return await axios

    .post("http://localhost/api/election/election", data)

    .then((res) => res)

    .catch((e) => e.response);

}

async function CreateVote(data: UsersInterface) {

  return await axios

    .post("http://localhost/api/vote/vote", data)

    .then((res) => res)

    .catch((e) => e.response);

}

async function GetVotes() {
  return await axios
    .get("http://localhost/api/vote/votes")
    .then((res) => res)
    .catch((e) => e.response);
}



export {

  SignIn,

  GetGender,

  GetUsers,

  GetUsersById,

  UpdateUsersById,

  DeleteUsersById,

  CreateUser,

  GetCandidates,
  GetCandidateById,
  UpdateCandidateById,
  DeleteCandidateById,
  CreateCandidate,

  GetElections,
  GetElectionById,
  UpdateElectionById,
  DeleteElectionById,
  CreateElection,

  CreateVote,
  GetVotes,
};