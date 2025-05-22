import { useEffect, useState } from "react";
import { Card, Button, Modal, message, Row, Col } from "antd";
import { useNavigate, useParams } from "react-router-dom";
import { GetCandidates, GetElectionById, CreateVote } from "../../services/https";
import { CandidateInterface } from "../../interfaces/ICandidate";
import { ElectionInterface } from "../../interfaces/IElection";

function VotePage() {
  const navigate = useNavigate();
  const { id } = useParams(); // electionId
  const [candidates, setCandidates] = useState<CandidateInterface[]>([]);
  const [election, setElection] = useState<ElectionInterface | null>(null);
  const [messageApi, contextHolder] = message.useMessage();
  const userId = localStorage.getItem("id");

  useEffect(() => {
    if (id) {
      fetchElection(id);
      fetchCandidates(id);
    }
  }, [id]);

  const fetchElection = async (electionId: string) => {
    const res = await GetElectionById(electionId);
    if (res.status === 200) {
      setElection(res.data);
    } else {
      messageApi.error("ไม่สามารถโหลดข้อมูลการเลือกตั้ง");
    }
  };

  const fetchCandidates = async (electionId: string) => {
    const res = await GetCandidates();
    if (res.status === 200) {
      const filtered = res.data.filter(
        (candidate: CandidateInterface) => String(candidate.election_id) === electionId
      );
      setCandidates(filtered);
    } else {
      messageApi.error("ไม่สามารถโหลดข้อมูลผู้สมัครได้");
    }
  };

  const confirmVote = (candidate: CandidateInterface) => {
    Modal.confirm({
      title: "ยืนยันการโหวต",
      content: `คุณต้องการโหวตให้ ${candidate.name} ใช่หรือไม่?`,
      okText: "ใช่, โหวตเลย",
      cancelText: "ยกเลิก",
      onOk: async () => {
        const res = await CreateVote({
          user_id: Number(userId),
          candidate_id: candidate.id,
          election_id: candidate.election_id,
          timestamp: new Date().toISOString(), // ISO 8601 format
        });

        if (res.status === 201) {
          messageApi.success("โหวตสำเร็จ");
          navigate("/elections");
        } else {
          messageApi.error(res.data?.error || "เกิดข้อผิดพลาดในการโหวต");
        }
      },
    });
  };

  return (
    <>
      {contextHolder}
      <h2>{election?.title || "กำลังโหลด..."}</h2>
      <p>{election?.description}</p>

      <Row gutter={[16, 16]}>
        {candidates.map((candidate) => (
          <Col key={candidate.ID} xs={24} sm={12} md={8}>
            <Card title={`${candidate.name}`}>
              <Button type="primary" onClick={() => confirmVote(candidate)}>
                โหวต
              </Button>
            </Card>
          </Col>
        ))}
      </Row>
    </>
  );
}

export default VotePage;
