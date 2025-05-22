import { useEffect, useState } from "react";
import { Card, Col, Row, message } from "antd";
import { GetCandidates } from "../../services/https";
import { useNavigate } from "react-router-dom";
import './index.css'; // เพิ่มบรรทัดนี้

interface Candidate {
  ID?: number;
  name?: string;
}

function Candidates() {
  const [candidates, setCandidates] = useState<Candidate[]>([]);
  const [messageApi, contextHolder] = message.useMessage();
  const navigate = useNavigate();

  const getCandidates = async () => {
    const res = await GetCandidates();
    if (res.status === 200) {
      setCandidates(res.data);
    } else {
      messageApi.error("ไม่สามารถโหลดข้อมูลผู้สมัครได้");
    }
  };

  useEffect(() => {
    getCandidates();
  }, []);

  return (
    <div className="page-container">
      {contextHolder}

      <Row justify="space-between" align="middle">
        <Col>
          <h2 className="page-title">รายชื่อผู้สมัคร</h2>
        </Col>
      </Row>

      <Row gutter={[16, 16]} className="card-container">
        {candidates.map((candidate) => (
          <Col xs={24} sm={12} md={8} lg={6} key={candidate.ID}>
            <Card title={candidate.name} bordered={true} />
          </Col>
        ))}
      </Row>
    </div>
  );
}

export default Candidates;
