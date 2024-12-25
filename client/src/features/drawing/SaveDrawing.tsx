import { Alert, Button, CircularProgress, Dialog, DialogActions, DialogContent, DialogTitle, TextField } from "@mui/material";
import { useAppDispatch, useAppSelector } from "../../app/hooks";
import { useState } from "react";
import { saveDrawingContent, selectCurrentDrawingContent, selectSavedDrawing, selectSaveDrawingStatus } from "./drawingSlice";

import style from "./Drawing.module.css";

interface SaveDrawingDialogProps {
  readonly open: boolean;
  readonly onClose: () => void;
}

export const SaveDrawingDialog = ({ open, onClose }: SaveDrawingDialogProps) => {

  const savedDrawing = useAppSelector(selectSavedDrawing);
  const currentContent = useAppSelector(selectCurrentDrawingContent);
  const savingStatus = useAppSelector(selectSaveDrawingStatus);

  const [selectedTitle, setSelectedTitle] = useState<string>("");

  const dispatch = useAppDispatch();

  const handleOk = () => {
    dispatch(saveDrawingContent({ title: selectedTitle || savedDrawing.title, content: currentContent }));
    onClose();
  };

  return (
    <Dialog
      open={open}
      onClose={onClose}
    >
      <DialogTitle>Save drawing</DialogTitle>
      <DialogContent>{
        <div className={style.openSaveDrawingDialogContent}>{
          savingStatus === "loading"
            ? <div className={style.fetchInProgress}>
                <CircularProgress />
              </div>
            : savingStatus === "failed"
              ? <Alert severity="error">Failed to save drawing</Alert>
              : <TitleSelector title={selectedTitle || savedDrawing.title} onChange={title => setSelectedTitle(title)}/>
        }</div>
      }
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Cancel</Button>
        <Button onClick={handleOk} autoFocus>OK</Button>
      </DialogActions>
    </Dialog>
  );
};

interface TitleSelectorProps {
  readonly title: string;
  readonly onChange: (selectedTitle: string) => void;
}

const TitleSelector = ({ title, onChange }:  TitleSelectorProps) => {
  return (
    <TextField label="Title" size="small" sx={{ width: "100%" }} onChange={event => onChange(event.target.value)} value={title}/>
  );
};