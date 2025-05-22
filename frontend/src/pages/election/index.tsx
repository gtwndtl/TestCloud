import { useEffect, useState } from "react";
import { Card, Button, Row, Col, message } from "antd";
import { useNavigate } from "react-router-dom";
import dayjs from "dayjs";

import { GetElections, GetVotes } from "../../services/https";
import { ElectionInterface } from "../../interfaces/IElection";
import { VoteInterface } from "../../interfaces/IVote";

import "./index.css"; // เพิ่มบรรทัดนี้

function Elections() {
  const [elections, setElections] = useState<ElectionInterface[]>([]);
  const [votes, setVotes] = useState<VoteInterface[]>([]);
  const [votedElectionIds, setVotedElectionIds] = useState<number[]>([]);
  const [messageApi, contextHolder] = message.useMessage();
  const navigate = useNavigate();

  const userId = localStorage.getItem("id");

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    const [electionRes, voteRes] = await Promise.all([
      GetElections(),
      GetVotes()
    ]);

    if (electionRes.status === 200) {
      setElections(electionRes.data);
    } else {
      messageApi.error("ไม่สามารถโหลดข้อมูลการเลือกตั้งได้");
    }

    if (voteRes.status === 200) {
      setVotes(voteRes.data);

      const userVotes = voteRes.data.filter(
        (vote: VoteInterface) => String(vote.user_id) === String(userId)
      );
      const votedIds = userVotes.map((v: VoteInterface) => v.election_id);
      setVotedElectionIds(votedIds);
    } else {
      messageApi.error("ไม่สามารถโหลดข้อมูลการโหวตได้");
    }
  };

  const goToElectionDetail = (id: number) => {
    navigate(`/election/${id}`);
  };

  const getLeadingCandidate = (electionId: number) => {
    const electionVotes = votes.filter(v => v.election_id === electionId);

    const voteCountMap: Record<number, number> = {};
    electionVotes.forEach(v => {
      voteCountMap[v.candidate_id] = (voteCountMap[v.candidate_id] || 0) + 1;
    });

    const sorted = Object.entries(voteCountMap).sort((a, b) => b[1] - a[1]);
    if (sorted.length === 0) return null;

    const [candidateId, count] = sorted[0];
    return { candidateId: Number(candidateId), count };
  };

  return (
    <div className="page-container">
      {contextHolder}
      <h2 className="page-title">รายการเลือกตั้ง</h2>
      <Row gutter={[16, 16]} className="card-wrapper">
        {elections.map((election) => {
          const hasVoted = votedElectionIds.includes(election.ID);
          const leader = getLeadingCandidate(election.ID);

          return (
            <Col key={election.ID} xs={24} sm={12} md={8}>
              <Card
                title={election.title}
                bordered
                style={{
                  backgroundColor: hasVoted ? "#f5f5f5" : "#ffffff",
                  borderColor: hasVoted ? "#d9d9d9" : "#1890ff"
                }}
                headStyle={{
                  backgroundColor: hasVoted ? "#d9d9d9" : "#1890ff",
                  color: "#fff"
                }}
                extra={<span style={{ color: "#fff" }}>{election.status}</span>}
              >
                <p>{election.description}</p>
                <p className="vote-info">
                  เริ่ม: {dayjs(election.start_time).format("DD/MM/YYYY HH:mm")}<br />
                  สิ้นสุด: {dayjs(election.end_time).format("DD/MM/YYYY HH:mm")}
                </p>

                {leader && (
                  <p className="leader-info">
                    ผู้นำ: หมายเลข {leader.candidateId} ({leader.count} คะแนน)
                  </p>
                )}

                <Button
                  type="primary"
                  block
                  disabled={hasVoted}
                  onClick={() => goToElectionDetail(election.ID)}
                  style={{ marginTop: 12 }}
                >
                  {hasVoted ? "คุณได้โหวตไปแล้ว" : "ดูรายละเอียด"}
                </Button>
              </Card>
            </Col>
          );
        })}
      </Row>
    </div>
  );
}

export default Elections;
