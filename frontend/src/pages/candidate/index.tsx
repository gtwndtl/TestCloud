import { useEffect, useState } from "react";
import { Card, Col, Row, message, Button } from "antd";
import { GetCandidates } from "../../services/https"; // สมมติว่ามีฟังก์ชันนี้
import { Link, useNavigate } from "react-router-dom";

interface Candidate {
  ID?: number;    // ตัวใหญ่ตาม json
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
    <>
      {contextHolder}

      <Row justify="space-between" align="middle" style={{ marginBottom: 16 }}>
        <Col>
          <h2>รายชื่อผู้สมัคร</h2>
        </Col>
      </Row>

      <Row gutter={[16, 16]}>
        {candidates.map((candidate) => (
          <Col xs={24} sm={12} md={8} lg={6} key={candidate.ID}>
            <Card title={candidate.name} bordered={true}>
            </Card>
          </Col>
        ))}
      </Row>
    </>
  );
}

export default Candidates;
