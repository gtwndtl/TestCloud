import {
  Space,
  Button,
  Col,
  Row,
  Divider,
  Form,
  Input,
  Card,
  message,
} from "antd";
import { PlusOutlined } from "@ant-design/icons";
import { } from "react";
import { useNavigate, Link } from "react-router-dom";
import { CreateCandidate } from "../../../services/https";

function CandidateCreate() {
  const navigate = useNavigate();
  const [messageApi, contextHolder] = message.useMessage();



  const onFinish = async (values: any) => {
    // ส่งข้อมูลไปสร้าง candidate
    const res = await CreateCandidate(values);
    if (res.status === 201) {
      messageApi.open({ type: "success", content: "เพิ่มผู้สมัครเรียบร้อย" });
      setTimeout(() => navigate("/candidate"), 2000);
    } else {
      messageApi.open({ type: "error", content: res.data.error || "เกิดข้อผิดพลาด" });
    }
  };

  return (
    <div>
      {contextHolder}
      <Card>
        <h2>เพิ่มข้อมูลผู้สมัคร</h2>
        <Divider />
        <Form name="candidate_create" layout="vertical" onFinish={onFinish} autoComplete="off">
          <Row gutter={[16, 0]}>
            <Col xs={24} sm={24} md={24} lg={12}>
              <Form.Item
                label="ชื่อผู้สมัคร"
                name="name"
                rules={[{ required: true, message: "กรุณากรอกชื่อผู้สมัคร!" }]}
              >
                <Input />
              </Form.Item>
            </Col>
          </Row>

          <Row justify="end">
            <Col style={{ marginTop: "40px" }}>
              <Form.Item>
                <Space>
                  <Link to="/candidate">
                    <Button htmlType="button" style={{ marginRight: "10px" }}>
                      ยกเลิก
                    </Button>
                  </Link>
                  <Button type="primary" htmlType="submit" icon={<PlusOutlined />}>
                    ยืนยัน
                  </Button>
                </Space>
              </Form.Item>
            </Col>
          </Row>
        </Form>
      </Card>
    </div>
  );
}

export default CandidateCreate;
